import {
  Page,
  Font,
  Text,
  View,
  Document,
  StyleSheet,
  Image,
  PDFViewer,
  Link,
} from '@react-pdf/renderer';
import matahariRegular from '../../fonts/Matahari-400Regular.ttf';
import montserratSemibold from '../../fonts/Montserrat-SemiBold.ttf';
import { useLocation } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { TransactionInfo } from '../../types/transactionInfo';
import { getAssetCode, getExplorerBaseUrl } from '../../utils/utilities';
import { truncateAddress } from '../../utils/truncateValues';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';

Font.register({
  family: 'Matahari',
  fonts: [
    { src: matahariRegular, fontWeight: 'normal' },
    { src: montserratSemibold, fontWeight: 'bold' },
  ],
});

// Font.registerHyphenationCallback((word) => {
//   return [word];
// });

// Create styles
const styles = StyleSheet.create({
  page: {
    flexDirection: 'row',
    padding: 60,
  },
  section: {
    margin: 10,
    padding: 10,
    flexGrow: 1,
  },
  header: {
    marginTop: '40px',
    marginBottom: '5px',
    color: '#004988',
    fontWeight: 'bold',
    fontFamily: 'Matahari',
    fontSize: '20pt',
    textAlign: 'center',
  },
  generated: {
    color: '#004988',
    fontFamily: 'Matahari',
    fontSize: '13pt',
    textAlign: 'center',
  },
  details: {
    marginTop: '40px',
    backgroundColor: '#e8ecf4',
    borderRadius: '20px',
    paddingVertical: '30px',
    paddingHorizontal: '20px',
  },
  itemHeader: {
    color: '#004988',
    fontFamily: 'Matahari',
    fontSize: '14pt',
    fontWeight: 'bold',
  },
  itemBody: {
    color: '#004988',
    fontFamily: 'Matahari',
    fontSize: '15pt',
    marginTop: '6pt',
    marginBottom: '4pt',
  },
  itemBody2: {
    color: '#004988',
    fontFamily: 'Matahari',
    fontSize: '12pt',
    marginBottom: '4pt',
  },
  transactionId: {
    color: '#004988',
    fontFamily: 'Matahari',
    fontSize: '10pt',
    marginTop: '6pt',
    marginBottom: '4pt',
    textOverflow: 'ellipsis',
    flex: 1,
  },
  link: {
    maxWidth: '400pt',
    display: 'flex',
    flexWrap: 'wrap',
    height: '10px',
    alignItems: 'center',
    color: '#004988',
    fontFamily: 'Matahari',
    fontSize: '10pt',
    marginTop: '6pt',
    marginBottom: '4pt',
  },
});

// Create Document Component
export default function SendAssetReceipt() {
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const appState = useSelector((state: RootState) => state.appState!);
  const [info, setInfo] = useState<TransactionInfo>();
  useEffect(() => {
    const b64Info = queryParams.get('q');
    setInfo(JSON.parse(atob(b64Info!)) as TransactionInfo);
  }, []);

  return (
    info && (
      <PDFViewer className="w-screen h-screen">
        <Document>
          <Page size="A4" style={styles.page}>
            <View style={styles.section}>
              <Image
                style={{ width: '350px', alignSelf: 'center' }}
                src="http://localhost:3000/images/trovoHorizontalLogo.png"
              />
              <Text style={styles.header}>Payment Details</Text>
              <Text style={styles.generated}>
                Generated from Trovo App on November 06, 2024 06:20 am
              </Text>
              <View style={styles.details}>
                <View
                  style={{
                    borderBottom: '1px',
                    borderBottomColor: '#d2d6dd',
                  }}
                >
                  <Text style={styles.itemHeader}>Sent from</Text>
                  <Text style={styles.itemBody}>{info!.from!}</Text>
                  <Text style={styles.itemBody2}>{info?.fromAddress}</Text>
                </View>
                <View
                  style={{
                    borderBottom: '1px',
                    borderBottomColor: '#d2d6dd',
                    marginTop: '15pt',
                  }}
                >
                  <Text style={styles.itemHeader}>To</Text>
                  <Text style={styles.itemBody}>{info?.to}</Text>
                  {info?.toAddress && (
                    <Text style={styles.itemBody2}>
                      {truncateAddress(info.toAddress)}
                    </Text>
                  )}
                </View>
                <View
                  style={{
                    borderBottom: '1px',
                    borderBottomColor: '#d2d6dd',
                    marginTop: '15pt',
                  }}
                >
                  <Text style={styles.itemHeader}>Amount</Text>
                  <Text style={styles.itemBody}>
                    {`${info?.amount} ${getAssetCode(info?.assetCode ?? '')}`}
                  </Text>
                </View>
                {info?.memo && (
                  <View
                    style={{
                      borderBottom: '1px',
                      borderBottomColor: '#d2d6dd',
                      marginTop: '15pt',
                    }}
                  >
                    <Text style={styles.itemHeader}>Memo</Text>
                    <Text style={styles.itemBody}>{info?.memo}</Text>
                  </View>
                )}
                <View
                  style={{
                    borderBottom: '1px',
                    borderBottomColor: '#d2d6dd',
                    marginTop: '15pt',
                  }}
                >
                  <Text style={styles.itemHeader}>
                    Blockchain Proof (Transaction ID)
                  </Text>
                  <Link
                    style={styles.link}
                    src={`${getExplorerBaseUrl(appState.walletMode)}${
                      info.transactionId
                    }`}
                  >
                    {info?.transactionId}
                  </Link>
                </View>
                <View
                  style={{
                    marginTop: '15pt',
                  }}
                >
                  <Text style={styles.itemHeader}>Date</Text>
                  <Text style={styles.itemBody}>
                    {new Date(info?.transactionDate!).toLocaleString()}
                  </Text>
                </View>
              </View>
            </View>
          </Page>
        </Document>
      </PDFViewer>
    )
  );
}
